package {{.Package}};

import static org.assertj.core.api.Assertions.assertThat;

import dev.galasa.Test;
import dev.galasa.core.manager.CoreManager;
import dev.galasa.core.manager.ICoreManager;
import dev.galasa.zos.IZosImage;
import dev.galasa.zos.ZosImage;
import dev.galasa.zos3270.ITerminal;
import dev.galasa.zos3270.Zos3270Terminal;

@Test
public class SampleTest {
	
	@ZosImage(imageTag = "PRIMARY")
	public IZosImage zosImage;
	
	@Zos3270Terminal(imageTag = "PRIMARY")
	public ITerminal zos3270Terminal;
	
	@CoreManager
	public ICoreManager core;
	
	@Test
	public void simpleSampleTest() {
		assertThat(zosImage).isNotNull();
		assertThat(zos3270Terminal).isNotNull();
	}
	
	@Test
	public void simpleLogonScreenTest() throws Exception {		
		assertThat(zos3270Terminal.retrieveScreen()).contains("VAMP");
	}

}
